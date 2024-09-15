library IEEE;

use IEEE.STD_LOGIC_1164.all;

entity AdderSubtractor is

    generic (
        N : integer := 5
        );
    port (
        A        : in  std_logic_vector (N-1 downto 0);
        B        : in  std_logic_vector (N-1 downto 0);
        nAdd_Sub : in  std_logic;
        Sum      : out std_logic_vector (N-1 downto 0);
        Carry    : out std_logic_vector (N-1 downto 0)
        );

end AdderSubtractor;

architecture Structural of AdderSubtractor is

    component NBitInverter
        port (
            Input  : in  std_logic_vector (N-1 downto 0);
            Output : out std_logic_vector (N-1 downto 0)
            );
    end component;

    component mux2t1_N
        generic (
            N : integer := 4
            );
        port (
            i_D0 : in  std_logic_vector (N-1 downto 0);
            i_D1 : in  std_logic_vector (N-1 downto 0);
            i_S  : in  std_logic;
            o_O  : out std_logic_vector (N-1 downto 0)
            );
    end component;

    component NBitAdder
        port (
            A        : in  std_logic_vector (N-1 downto 0);
            B        : in  std_logic_vector (N-1 downto 0);
            Sum      : out std_logic_vector (N-1 downto 0);
            CarryOut : out std_logic_vector (N-1 downto 0)
            );
    end component;

    signal s_inverted : std_logic_vector (N-1 downto 0);
    signal s_muxed    : std_logic_vector (N-1 downto 0);

begin

    Inv : NBitInverter
        port map (
            Input  => B,
            Output => s_inverted
            );

    Mux : mux2t1_N
        port map (
            i_D0 => B,
            i_D1 => s_inverted,
            i_S  => nAdd_Sub,
            o_O  => s_muxed
            );

    -- Instantiate the N-bit adder
    Adder : NBitAdder
        port map (
            A        => A,
            B        => s_muxed,
            Sum      => Sum,
            CarryOut => Carry
            );

end Structural;
