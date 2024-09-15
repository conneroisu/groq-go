library ieee;


use IEEE.STD_LOGIC_1164.all;
use IEEE.STD_LOGIC_ARITH.all;
use IEEE.STD_LOGIC_UNSIGNED.all;
use IEEE.numeric_std.all;

entity tb_Adder is end tb_Adder;

architecture behavior of tb_Adder is
    component Adder
        generic (
            N : integer := 4
            );
        port (
            A        : in  std_logic_vector (N-1 downto 0);
            B        : in  std_logic_vector (N-1 downto 0);
            nAdd_Sub : in  std_logic;
            Sum      : out std_logic_vector (N-1 downto 0);
            Carry    : out std_logic_vector (N-1 downto 0)
            );
    end component;

    signal s_A       : std_logic_vector (3 downto 0) := (others => '0');
    signal s_B       : std_logic_vector (3 downto 0) := (others => '0');
    signal s_nAddSub : std_logic                     := '0';
    signal s_Carry   : std_logic_vector (3-1 downto 0);

    signal s_Sum : std_logic_vector (3 downto 0);

begin
    DUT0 : Adder generic map (N => 4)
        port map (
            A        => s_A,
            B        => s_B,
            nAdd_Sub => s_nAddSub,
            Sum      => s_Sum,
            Carry    => s_Carry
            );

    s_A <= "0011"; s_B <= "0101"; s_nAddSub <= '0';
    process
    begin
        wait for 10 ns;

        s_A <= "0110"; s_B <= "0011"; s_nAddSub <= '1';
        wait for 10 ns;

        s_A <= "0000"; s_B <= "0000"; s_nAddSub <= '0';  -- Add zero
        wait for 10 ns;
        s_A <= "0000"; s_B <= "0000"; s_nAddSub <= '1';  -- Subtract zero
        wait for 10 ns;
        s_A <= "1111"; s_B <= "1111"; s_nAddSub <= '0';  -- Add max
        wait for 10 ns;
        s_A <= "1111"; s_B <= "1111"; s_nAddSub <= '1';  -- Subtract max
        wait for 10 ns;

        s_A <= "1111"; s_B <= "0001"; s_nAddSub <= '0';  -- Add overflow
        wait for 10 ns;
        s_A <= "0000"; s_B <= "0001"; s_nAddSub <= '1';  -- Subtract overflow
        wait for 10 ns;

        assert false report "Testbench completed" severity note;
        wait;
    end process;
end;
