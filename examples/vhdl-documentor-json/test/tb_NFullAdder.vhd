library IEEE;
use IEEE.STD_LOGIC_1164.all;
use IEEE.NUMERIC_STD.all;

entity tb_NFullAdder is
-- Empty entity as this is a test bench
end tb_NFullAdder;

architecture Behavioral of tb_NFullAdder is
    component NBitAdder
        Generic (N : integer := 4);
        Port (
            A : in STD_LOGIC_VECTOR(N-1 downto 0);
            B : in STD_LOGIC_VECTOR(N-1 downto 0);
            Sum : out STD_LOGIC_VECTOR(N-1 downto 0);
            CarryOut : out STD_LOGIC
        );
    end component;

    signal s_A : STD_LOGIC_VECTOR(3 downto 0) := (others => '0');
    signal s_B : STD_LOGIC_VECTOR(3 downto 0) := (others => '0');

    signal s_Sum : STD_LOGIC_VECTOR(3 downto 0);
    signal s_CarryOut : STD_LOGIC;

    signal clk : STD_LOGIC := '0';

begin
    DUT0: NBitAdder 
        port map (
            A => s_A,
            B => s_B,
            Sum => s_Sum,
            CarryOut => s_CarryOut
        );

    -- Clock process
    clk_process : process
    begin
        clk <= '0';
        wait for 5 ns;
        clk <= '1';
        wait for 5 ns;
    end process;

    -- Test process
    stim_proc: process
    begin
        -- Testing all combinations
        for i in 0 to 15 loop
            for j in 0 to 15 loop
                s_A <= std_logic_vector(to_unsigned(i, 4));
                s_B <= std_logic_vector(to_unsigned(j, 4));
                wait for 10 ns;  -- Wait for one clock cycle
            end loop;
        end loop;

        -- End simulation
        wait;
    end process;

end Behavioral;
